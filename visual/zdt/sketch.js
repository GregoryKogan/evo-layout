let data = [];
let currentGen = 0;
let loaded = false;

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  // Load the JSONL log file (each line is a JSON object)
  loadStrings('ZDT1_NSGA2.jsonl', function(lines) {
    // Parse each line as JSON
    data = lines.filter(l => l.trim().length > 0).map(l => JSON.parse(l));
    data.shift()
    loaded = true;
    frameRate(5);
    saveGif("nsga2-zdt1.gif", 20);
  });
}

function draw() {
  if (!loaded) return;
  background(18);
  
  // Get the current log entry (if available)
  let entry = data[currentGen];
  if (entry) { 
    // For visualization, assume objectives f1 and f2 are in [0,1] (or scale appropriately)
    // Draw axes
    stroke(255);
    line(50, height - 50, width - 50, height - 50); // x-axis
    line(50, height - 50, 50, 50); // y-axis
    // Label axes
    noStroke();
    fill(255);
    textSize(16);
    text("Generation: " + entry.generation, 20, 30)
    text("f1", width - 40, height - 30);
    text("f2", 20, 60);
    
    // Plot each point from the Pareto front
    for (let point of entry.pareto_front) {
      let x = map(point.f1, 0, 1, 50, width - 50);
      let y = map(point.f2, 0, 1, height - 50, 50); // invert y-axis
      fill(0, 200, 255);
      noStroke();
      ellipse(x, y, 10, 10);
    }
    
    // Advance to next generation on each frame, stopping at the end.
    currentGen += (currentGen + 1 < data.length)
  }
}
