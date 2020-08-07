const body = document.getElementsByTagName("body")[0];
const search = document.getElementById("search");
let results = document.getElementById("results");

search.onkeyup = (event) => {
    clearResults();

    if (!search.value) {
        return;
    }

    getMatches(search.value)
        .then((matches) => drawResults(results, matches));
};

async function getMatches(term) {
    try {
        const resp = await fetch(`/autocomplete?term=${term}`);
        const words = (await resp.text()).split("\n").filter((w) => w.length);
        return words;
    } catch (err) {
        console.error(err);
        return [];
    }
}

function clearResults() {
    results.remove();
    results = document.createElement("ul");
    results.setAttribute("id", "results");
    body.appendChild(results);
}

function drawResults(results, matches) {
    matches.forEach((match) => {
        const li = document.createElement("li");
        li.innerText = match;
        results.appendChild(li);
    });
}
