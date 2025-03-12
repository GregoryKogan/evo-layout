let solutions;
let loaded = false;
let saved = false;

let colors = {};

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
        solutions = parsed.slice(2);
      } else {
        solutions = parsed.slice(1).map((x) => {
          return x.solution.order;
        });
      }
      loaded = true;
      assignColors(solutions[solutions.length - 1]);
    });
  });
  frameRate(3);
  // saveGif("planar-graph.gif", 8);
}

function draw() {
  if (!loaded) return;

  background(18);

  // drawLines();
  drawCircles();
  // drawPixels();
}

function drawLines() {
  const solutionsNum = solutions.length;
  const citiesNum = solutions[0].length;
  for (let i = 1; i < solutionsNum; i++) {
    for (let j = 0; j < citiesNum; j++) {
      const [x1, y1] = toCoords(solutionsNum, citiesNum, i, j);
      prevIndex = solutions[i - 1].indexOf(solutions[i][j]);
      const [x2, y2] = toCoords(solutionsNum, citiesNum, i - 1, prevIndex);
      stroke(colors[solutions[i][j]]);
      strokeWeight(2);
      line(x1, y1, x2, y2);
    }
  }

  if (!saved) {
    // save("convergence-lines.png");
    saved = true;
  }
}

function drawCircles() {
  const solutionsNum = solutions.length;
  const citiesNum = solutions[0].length;
  const r = min(width / citiesNum, height / solutionsNum);
  for (let i = 0; i < solutionsNum; i++) {
    for (let j = 0; j < citiesNum; j++) {
      const [x, y] = toCoords(solutionsNum, citiesNum, i, j);
      fill(colors[solutions[i][j]]);
      noStroke();
      circle(x, y, r);
    }
  }

  if (!saved) {
    // save("convergence-circles.png");
    saved = true;
  }
}

function drawPixels() {
  const solutionsNum = solutions.length;
  const citiesNum = solutions[0].length;
  const pw = width / citiesNum;
  const ph = height / solutionsNum;
  for (let i = 0; i < solutionsNum; i++) {
    for (let j = 0; j < citiesNum; j++) {
      const [x, y] = toCoords(solutionsNum, citiesNum, i, j);
      fill(colors[solutions[i][j]]);
      noStroke();
      rect(x - pw / 2, y - ph / 2, pw, ph);
    }
  }

  if (!saved) {
    // save("convergence-pixels.png");
    saved = true;
  }
}

function assignColors(lastSolution) {
  const citiesNum = lastSolution.length;
  colorMode(HSB);
  angleMode(DEGREES);
  const step = 360.0 / citiesNum;
  let angle = 0;
  for (city of lastSolution) {
    colors[city] = color(angle, 85, 90);
    angle += step;
  }
}

function toCoords(solutionsNum, citiesNum, i, j) {
  x = (width / citiesNum) * j + width / citiesNum / 2;
  y = (height / solutionsNum) * i + height / solutionsNum / 2;
  return [x, y];
}
