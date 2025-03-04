let problem;
let algoSolution;
let solutions;
let loaded = false;
let done = false;

let generation = 0;

function setup() {
  createCanvas(min(windowWidth, windowHeight), min(windowWidth, windowHeight));
  textSize(32);
  fetch("TSP.jsonl").then((response) => {
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
      if (parsed[1].took !== undefined) {
        algoSolution = parsed[1].solution;
        solutions = parsed.slice(2);
      } else {
        solutions = parsed.slice(1);
      }
      loaded = true;
    });
  });
  frameRate(3);
  // saveGif("planar-graph.gif", 8);
}

function draw() {
  if (!loaded || done) return;

  if (generation >= solutions.length) {
    done = true;
    console.log("Animation frames", frameCount);
    return;
  }

  const solution = solutions[generation].solution;
  generation++;

  background(18);

  document.querySelector("#length").innerHTML = `Length: ${round(
    1 / solution.fitness
  )}`;

  for (let city of problem.cities) {
    fill(255);
    const v = toScreenCoord(city.lat, city.lon);
    circle(v.x, v.y, 10);
  }

  if (algoSolution) drawSolution(algoSolution, "#00ff00");
  drawSolution(solution, 255);
}

function toScreenCoord(x, y) {
  x /= 100;
  y /= 100;
  return createVector(50 + x * (width - 100), 50 + y * (height - 100));
}

function drawSolution(solution, color) {
  let curCity = problem.cities[0];
  for (let cityInd of solution.order) {
    const city = problem.cities[cityInd];
    stroke(color);
    strokeWeight(2);
    line(
      toScreenCoord(curCity.lat, curCity.lon).x,
      toScreenCoord(curCity.lat, curCity.lon).y,
      toScreenCoord(city.lat, city.lon).x,
      toScreenCoord(city.lat, city.lon).y
    );
    curCity = city;
  }
  line(
    toScreenCoord(curCity.lat, curCity.lon).x,
    toScreenCoord(curCity.lat, curCity.lon).y,
    toScreenCoord(problem.cities[0].lat, problem.cities[0].lon).x,
    toScreenCoord(problem.cities[0].lat, problem.cities[0].lon).y
  );
}
