let problem;
let solutions;
let loaded = false;
let done = false;

let generation = 0;

function setup() {
  createCanvas(min(windowWidth, windowHeight), min(windowWidth, windowHeight));
  textSize(32);
  fetch("GraphPlane_SGA.jsonl").then((response) => {
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
  frameRate(5);
  // saveGif("planar-graph.gif", 6);
}

function draw() {
  if (!loaded || done) return;
  generation++;
  if (generation >= solutions.length) {
    done = true;
    console.log("Animation frames", frameCount);
    return;
  }

  background(18);

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
