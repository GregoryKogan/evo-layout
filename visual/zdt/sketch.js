let data = [];
let currentGen = 0;
let loaded = false;

// Define a zoom factor for screen margins.
const zoom = 0.95;

// Set the ZDT problem type to visualize ("ZDT1", "ZDT2", "ZDT3", "ZDT4", or "ZDT6").
const zdtType = "ZDT1";
const algorithm = "NSGA2";
const gifLenSeconds = 20;

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  // Load the JSONL log file (each line is a JSON object)
  loadStrings(zdtType + '_' + algorithm + '.jsonl', function(lines) {
    // Parse each line as JSON and remove the header if needed.
    data = lines.filter(l => l.trim().length > 0).map(l => JSON.parse(l));
    data.shift();
    loaded = true;
    frameRate(Math.ceil(data.length / gifLenSeconds));
    // Uncomment to save a GIF when running in an environment that supports it:
    // saveGif(zdtType + '-' + algorithm + '.gif', gifLenSeconds);
  });
}

function draw() {
  if (!loaded) return;
  background(18);
  
  // Draw axes
  stroke(255);
  line(30, height - 30, width * zoom, height - 30); // x-axis
  line(30, height - 30, 30, height * (1 - zoom));     // y-axis
  
  // Draw the true Pareto front using the selected ZDT type.
  drawTrueParetoFront();

  let entry = data[currentGen];
  if (entry) { 
    // Label axes and generation/iteration info.
    noStroke();
    fill(255);
    textSize(16);
    textAlign(RIGHT);
    text(zdtType + ' ' + algorithm, width - 20, 20)
    textAlign(LEFT);
    if (entry.generation) text("Generation: " + entry.generation, 20, 20);
    if (entry.iteration) text("Iteration: " + entry.iteration, 20, 20);
    text("f1", width * zoom - 10, height - 10);
    text("f2", 10, height * (1 - zoom) + 10);

    // Plot each point from the current Pareto front as blue circles.
    for (let point of entry.pareto_front) {
      drawPoint(point);
    }
    // If a best solution is logged separately, draw it with a red outline.
    if (entry.solution) {
      drawPoint(entry.solution.objectives, true);
    }
    
    // Advance to next generation on each frame, if available.
    if (currentGen + 1 < data.length) {
      currentGen++;
    }
  }
}

// Draws the true Pareto front for the selected ZDT problem.
function drawTrueParetoFront() {
  stroke(0, 255, 0);
  strokeWeight(2);
  noFill();
  
  if (zdtType === "ZDT3") {
    // The Pareto front for ZDT3 is discontinuous.
    // Use the five known intervals (from literature) and a fixed sample count.
    let intervals = [
      [0.0, 0.0830015349],
      [0.1822287280, 0.2577623634],
      [0.4093136748, 0.4538821041],
      [0.6183967944, 0.6525117038],
      [0.8233317983, 0.8518328654]
    ];
    for (let interval of intervals) {
      beginShape();
      let samples = 100; // fixed number of samples per interval
      for (let i = 0; i <= samples; i++) {
        let x = lerp(interval[0], interval[1], i / samples);
        // True f2 for ZDT3 (when g=1):
        let y = 1 - sqrt(x) - x * sin(10 * PI * x);
        let screenX = map(x, 0, 1, 30, width * zoom);
        let screenY = mapF2(y);
        vertex(screenX, screenY);
      }
      endShape();
    }
    return;
  }
  
  // For other ZDTs:
  if (zdtType === "ZDT6") {
    // For ZDT6, sample a parameter t from 0 to 1, compute the actual f1 value using its expression,
    // then define f2 = 1 - (f1)^2.
    beginShape();
    for (let t = 0; t <= 1; t += 0.0001) {
         let f1_val = 1 - exp(-4 * t) * pow(sin(6 * PI * t), 6);
         let f2_val = 1 - f1_val * f1_val;
         let screenX = map(f1_val, 0, 1, 30, width * zoom);
         let screenY = mapF2(f2_val);
         vertex(screenX, screenY);
    }
    endShape();
    return;
  }
  
  // For ZDT1, ZDT2, and ZDT4, assume the known front is given by a continuous relationship.
  beginShape();
  for (let x = 0; x <= 1; x += 0.01) {
    let y;
    if (zdtType === "ZDT1" || zdtType === "ZDT4") {
      y = 1 - sqrt(x);
    } else if (zdtType === "ZDT2") {
      y = 1 - x * x;
    }
    let screenX = map(x, 0, 1, 30, width * zoom);
    let screenY = mapF2(y);
    vertex(screenX, screenY);
  }
  endShape();
}

// Draws a point given an array [f1, f2].
// If best is true, the point is highlighted with a red outline.
function drawPoint(point, best = false) {
  let x = map(point[0], 0, 1, 30, width * zoom);
  let y = mapF2(point[1]);
  noStroke();
  fill(0, 200, 255);
  circle(x, y, 10);
  if (best) {
    fill(255, 0, 0);
    circle(x, y, 12);
  }
}

// Helper function for mapping f2 values.
// For ZDT3, we use preset bounds (e.g. from -0.3 to 1.0) because some f2 values may be negative.
function mapF2(f2) {
    if (zdtType === "ZDT3") {
        const minVal = 1 - sqrt(0.8518328654) - 0.8518328654 * sin(10 * PI * 0.8518328654)
        return map(f2, minVal, 1.0, height - 30, height * (1 - zoom));
    } else {
        return map(f2, 0, 1, height - 30, height * (1 - zoom));
    }
}
  