let problem;
let solutions;
let loaded = false;
let done = false;

let generation = 0;

function setup() {
  createCanvas(min(windowWidth, windowHeight), min(windowWidth, windowHeight));
  textSize(16);
  fill(255);
  stroke(255);
  textAlign(LEFT);
  fetch("PlanarGraphPlane_NSGA2+Force.jsonl").then((response) => {
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
  frameRate(25);
  // saveGif("gp-50-planar-Force->NSGA2.gif", 10);
}

function draw() {
  if (!loaded || done) return;
  generation++;
  if (solutions[generation] && solutions[generation].generation == null) {
    generation += 8;
    if (solutions[generation].solution.intersections == solutions[generation - 8].solution.intersections) generation += 32;
  } else if (solutions[generation] && solutions[generation].generation != null) {
    generation += 8;
  }
  if (generation >= solutions.length) {
    done = true;
    generation = solutions.length - 1;
    console.log("Animation frames", frameCount);
  }


  background(18);

  noStroke();
  if (solutions[generation].generation != null) text("Algorithm: NSGA2", 30, 50);
  else text("Algorithm: Fruchterman–Reingold", 30, 50);
  text("Intersections: " + solutions[generation].solution.intersections, 30, 30);

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

  for (let vertex of solutions[generation].solution.vertices) {
    fill(255);
    const v = toScreenCoord(vertex.x, vertex.y);
    circle(v.x, v.y, 10);
  }
}

function toScreenCoord(x, y) {
  return createVector(50 + x * (width - 100), 50 + y * (height - 100));
}
