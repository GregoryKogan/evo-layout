let data = [];
let currentGen = 0;
let generationIncrement = 1;
let loaded = false;

// Define a zoom factor for screen margins.
const zoom = 0.95;

const filename = 'GraphPlane_SPEA2'
const gifLenSeconds = 10;
const maxFrameRate = 30;
const generateGif = false;

const objTitles = ['intersections', 'min-angle', 'edge-length', 'dispersion']
const xInd = 2
const yInd = 0

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  loadStrings(filename + '.jsonl', function(lines) {
    data = lines.filter(l => l.trim().length > 0).map(l => JSON.parse(l));
    data.shift();
    loaded = true;
    fr = min(ceil(data.length / gifLenSeconds), maxFrameRate)
    frameRate(fr);
    generationIncrement = ceil(data.length / (fr * gifLenSeconds))
    if (generateGif) saveGif(filename + '_' + objTitles[xInd] + '_' + objTitles[yInd] + '.gif', gifLenSeconds);
  });
}

function draw() {
  if (!loaded) return;
  background(18);
  
  // Draw axes
  stroke(255);
  line(30, height - 30, width * zoom, height - 30); // x-axis
  line(30, height - 30, 30, height * (1 - zoom));   // y-axis

  let entry = data[currentGen];
  if (entry) { 
    // Label axes and generation/iteration info.
    noStroke();
    fill(255);
    textSize(16);
    textAlign(RIGHT);
    text(filename, width - 20, 20)
    textAlign(LEFT);
    if (entry.generation) text("Generation: " + entry.generation, 20, 20);
    if (entry.iteration) text("Iteration: " + entry.iteration, 20, 20);
    text(objTitles[yInd], 40, height * (1 - zoom) + 10);
    textAlign(RIGHT);
    text(objTitles[xInd], width * zoom - 10, height - 10);

    // Plot each point from the current Pareto front as blue circles.
    for (let point of entry.pareto_front) drawPoint(point);
    
    // If a best solution is logged separately, draw it with a red outline.
    if (entry.solution) drawPoint(entry.solution.objectives, true);
    
    // Advance to next generation on each frame, if available.
    if (currentGen + generationIncrement < data.length) currentGen += generationIncrement;
    else currentGen = data.length - 1;
  }
}

// Draws a point given an array [f1, f2].
// If best is true, the point is highlighted with a red outline.
function drawPoint(point, best = false) {
  let x = map(point[xInd], 0, 1, 30, width * zoom);
  let y = map(point[yInd], 0, 1, height - 30, height * (1 - zoom));
  noStroke();
  fill(0, 200, 255);
  circle(x, y, 10);
  if (best) {
    fill(255, 0, 0);
    circle(x, y, 12);
  }
}
