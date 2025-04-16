let data = [];
let currentGen = 0;
let loaded = false;

const zoom = 0.95

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  // Load the JSONL log file (each line is a JSON object)
  loadStrings('ZDT1_SSGA.jsonl', function(lines) {
    // Parse each line as JSON, remove the first log entry if needed.
    data = lines.filter(l => l.trim().length > 0).map(l => JSON.parse(l));
    data.shift();
    loaded = true;
    frameRate(5);
    // saveGif("nsga2-zdt1.gif", 20);
  });
}

function draw() {
  if (!loaded) return;
  background(18);
  
  // Draw axes
  stroke(255);
  line(30, height - 30, width * zoom, height - 30); // x-axis
  line(30, height - 30, 30, height * (1 - zoom)); // y-axis
  
  
  // Draw the correct (true) Pareto front in red.
  drawTrueParetoFront();

  let entry = data[currentGen];
  if (entry) { 
    // Label axes and generation info
    noStroke();
    fill(255);
    textSize(16);
    if (entry.generation) text("Generation: " + entry.generation, 20, 20);
    if (entry.iteration) text("Iteration: " + entry.iteration, 20, 20);
    text("f1", width * zoom - 10, height - 10);
    text("f2", 10, height * (1 - zoom) + 10);

    // Plot each point from the current Pareto front using a blue color.
    for (let point of entry.pareto_front) drawPoint(point);
    if (entry.solution) drawPoint(entry.solution.objectives, best=true);
    
    // Advance to next generation on each frame, stopping at the end.
    if (currentGen + 1 < data.length) {
      currentGen++;
    }
  }
}

// Draw the true Pareto front for ZDT1 in green
function drawTrueParetoFront() {
  stroke(0, 255, 0);
  strokeWeight(2);
  noFill();
  beginShape();
  // Draw points along the curve f2 = 1 - sqrt(f1)
  for (let x = 0; x <= 1; x += 0.01) {
    let y = 1 - sqrt(x);
    let screenX = map(x, 0, 1, 30,  width * zoom);
    let screenY = map(y, 0, 1, height - 30, height * (1 - zoom)); // invert y-axis
    vertex(screenX, screenY);
  }
  endShape();
}

function drawPoint(point, best=false) {
    let x = map(point[0], 0, 1, 30, width * zoom );
    let y = map(point[1], 0, 1, height - 30, height * (1 - zoom)); // invert y-axis
    noStroke();
    fill(0, 200, 255);
    circle(x, y, 10);
    if (best) {
        fill(255, 0, 0);
        circle(x, y, 12);
    }
}