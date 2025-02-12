let problem;
let solutions;
let loaded = false;

let generation = 0;

function setup() {
  createCanvas(windowWidth, windowHeight);
  background(18);
  fetch("results.json").then((response) => {
    response.json().then((data) => {
      problem = data.problem;
      solutions = data.bestSolutions;
      loaded = true;
    });
  });
  frameRate(10);
  saveGif("planar-graph.gif", 9);
}

function draw() {
  if (!loaded) return;
  background(18);

  generation++;
  if (generation >= solutions.length) generation = solutions.length - 1;
  else console.log("CUM");

  for (let vertex of solutions[generation].solution.vertices) {
    fill(255);
    const v = toScreenCoord(vertex.x, vertex.y);
    circle(v.x, v.y, 10);
  }

  for (let edge of problem.graph.edges) {
    let v1 = solutions[generation].solution.vertices[edge.from];
    let v2 = solutions[generation].solution.vertices[edge.to];
    stroke(255);
    strokeWeight(2);
    line(
      toScreenCoord(v1.x, v1.y).x,
      toScreenCoord(v1.x, v1.y).y,
      toScreenCoord(v2.x, v2.y).x,
      toScreenCoord(v2.x, v2.y).y
    );
  }
}

function toScreenCoord(x, y) {
  return createVector(50 + x * (width - 100), 50 + y * (height - 100));
}
