let problem;
let solutions;
let loaded = false;
let done = false;
let algorithm = "Fruchterman–Reingold";

let step = 0;

function setup() {
  createCanvas(min(windowWidth, windowHeight), min(windowWidth, windowHeight));
  textSize(16);
  fill(255);
  stroke(255);
  textAlign(LEFT);
  fetch("PlanarGraphPlane_FR-NSGA2.jsonl").then((response) => {
    response.text().then((data) => {
      const parsed = JSON.parse(
        "[" +
        data
          .split("\n")
          .filter((line) => line.length > 2)
          .join(",") +
        "]"
      );
      problem = parsed[0];
      solutions = parsed.slice(1);
      loaded = true;
    });
  });
  frameRate(60);
  // saveGif("gp-200-planar-FR-NSGA2.gif", 16);
}

function draw() {
  if (!loaded || done) return;

  curIntersections = solutions[step].solution.intersections;
  for (let fd = 0; fd < 8; ++fd) {
    step += (step < solutions.length - 1);
    if (solutions[step].step < solutions[step - 1].step) algorithm = "NSGA2";

    if (algorithm != "NSGA2" && solutions[step].solution.intersections < curIntersections && fd >= 0) break;
    if (algorithm == "NSGA2" && fd >= 2) break;
  }

  if (step >= solutions.length - 1) {
    done = true;
    step = solutions.length - 1;
    console.log("Animation frames", frameCount);
  }


  background(18);

  noStroke();
  // text("Algorithm: " + algorithm, 30, 50);
  // text("Intersections: " + solutions[step].solution.intersections, 30, 30);

  for (let edge of problem.graph.edges) {
    let v1 = solutions[step].solution.vertices[edge.from];
    let v2 = solutions[step].solution.vertices[edge.to];
    stroke(255);
    strokeWeight(1.4);
    line(
      toScreenCoord(v1.x, v1.y).x,
      toScreenCoord(v1.x, v1.y).y,
      toScreenCoord(v2.x, v2.y).x,
      toScreenCoord(v2.x, v2.y).y
    );
  }

  for (const vertex of solutions[step].solution.vertices) {
    fill(255);
    const v = toScreenCoord(vertex.x, vertex.y);
    circle(v.x, v.y, 7);
  }


}

function toScreenCoord(x, y) {
  return createVector(10 + x * (width - 20), 10 + y * (height - 20));
}

function keyPressed() {
  if (keyCode === ENTER) {
    saveCanvas();
  }
}