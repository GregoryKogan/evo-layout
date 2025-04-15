let data = [];
let currentGen = 0;
let loaded = false;

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  // Load the JSONL log file (each line is a JSON object)
  loadStrings('ZDT1_NSGA2.jsonl', function(lines) {
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
  line(50, height - 50, width - 50, height - 50); // x-axis
  line(50, height - 50, 50, 50); // y-axis
  
  
  // Draw the correct (true) Pareto front in red.
  drawTrueParetoFront();

  let entry = data[currentGen];
  if (entry) { 
    // Label axes and generation info
    noStroke();
    fill(255);
    textSize(16);
    text("Generation: " + entry.generation, 20, 30);
    text("f1", width - 40, height - 30);
    text("f2", 20, 60);

    // Plot each point from the current Pareto front using a blue color.
    for (let point of entry.pareto_front) {
      let x = map(point.f1, 0, 1, 50, width - 50);
      let y = map(point.f2, 0, 1, height - 50, 50); // invert y-axis
      fill(0, 200, 255);
      noStroke();
      ellipse(x, y, 10, 10);
    }
    
    // Advance to next generation on each frame, stopping at the end.
    if (currentGen + 1 < data.length) {
      currentGen++;
    }
  }
}

// Draw the true Pareto front for ZDT1 in red.
function drawTrueParetoFront() {
  stroke(0, 255, 0);
  strokeWeight(2);
  noFill();
  beginShape();
  // Draw points along the curve f2 = 1 - sqrt(f1)
  for (let x = 0; x <= 1; x += 0.01) {
    let y = 1 - sqrt(x);
    let screenX = map(x, 0, 1, 50, width - 50);
    let screenY = map(y, 0, 1, height - 50, 50); // invert y-axis
    vertex(screenX, screenY);
  }
  endShape();
}
