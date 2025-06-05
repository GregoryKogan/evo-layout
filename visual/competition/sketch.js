let data = [];
let currentGen = 2000;
let generationIncrement = 1;
let loaded = false;

// Define a zoom factor for screen margins.
const zoom = 0.99;

const filename = 'PlanarGraphPlane_FR-NSGA2';
const gifLenSeconds = 10;
const maxFrameRate = 30;
const generateGif = false;

const objTitles = ['intersections', 'dispersion', 'angle'];
const xInd = 1;
const yInd = 2;

function setup() {
  createCanvas(min(windowWidth, windowHeight), min(windowWidth, windowHeight));
  background(18);
  loadStrings(filename + '.jsonl', function (lines) {
    data = lines.filter(l => l.trim().length > 0).map(l => JSON.parse(l));
    data.shift();
    loaded = true;
    fr = min(ceil(data.length / gifLenSeconds), maxFrameRate);
    frameRate(fr);
    generationIncrement = ceil(data.length / (fr * gifLenSeconds));
    if (generateGif) saveGif(filename + '_' + objTitles[xInd] + '_' + objTitles[yInd] + '.gif', gifLenSeconds);
  });
}

function draw() {
  if (!loaded) return;
  background(255);

  // Draw axes
  stroke(0);
  strokeWeight(3);
  line(50, height - 30, width * zoom, height - 30); // x-axis
  line(50, height - 30, 50, height * (1 - zoom));   // y-axis

  // let entry = data[data.length - 1];
  let entry = data[3396];
  if (entry) {
    // Label axes and generation/iteration info.
    noStroke();
    fill(0);
    textSize(16);
    textAlign(RIGHT);
    text("FR-NSGA2", width - 20, 20);
    textAlign(LEFT);
    if (entry.generation) text("Generation: " + entry.generation, 20, 20);
    if (entry.iteration) text("Iteration: " + entry.iteration, 20, 20);
    text(objTitles[yInd], 60, height * (1 - zoom) + 10);
    text("0.800", 5, height * (1 - zoom) + 10);
    text("0.785", 5, height - 33);
    textAlign(RIGHT);
    text(objTitles[xInd], width * zoom - 5, height - 37);
    text("1.760", width * zoom - 5, height - 10);
    text("1.685", 90, height - 10);

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
  let x = map(point[xInd], 1.685, 1.760, 30, width * zoom);
  let y = map(point[yInd], 0.785, 0.8, height - 30, height * (1 - zoom));
  noStroke();
  fill("#d62729");
  circle(x, y, 10);
  // if (best) {
  //   fill(255, 0, 0);
  //   circle(x, y, 12);
  // }
}

function keyPressed() {
  if (keyCode === ENTER) {
    saveCanvas();
  }
}