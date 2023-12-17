import localforage from "localforage";
import { matchSorter } from "match-sorter";
import sortBy from "sort-by";

const globalData = new Map([
  [
    "1",
    {
      name: "OpenAI",
      description: "OpenAI Embeddings",
      data: {
        "2d": [
          [0.5, 0.6],
          [0.1, 0.2],
          [0.4, 0.1],
        ],
        "3d": [
          [0.4, 0.1, 0.3],
          [0.2, 0.1, 0.1],
          [0.1, 0.1, 0.2],
        ],
      },
    },
  ],
  [
    "2",
    {
      name: "VertexAI",
      description: "VertexAI Embeddings",
      data: {
        "2d": [
          [0.43, 0.77],
          [0.21, 0.33],
          [0.42, 0.78],
        ],
        "3d": [
          [0.234, 0.144, 0.453],
          [0.4442, 0.661, 0.51],
          [0.22, 0.11, 0.82],
        ],
      },
    },
  ],
]);

export async function getProviders(query) {
  await fakeNetwork(`getProviders:${query}`);
  let providers = [];
  for (const [id, obj] of globalData) {
    providers.push({ id: id, name: obj.name });
  }
  if (!providers) return [];
  if (query) {
    providers = matchSorter(providers, query, { keys: ["name"] });
  }
  return providers.sort(sortBy("name"));
}

export async function getProvider(id) {
  await fakeNetwork(`provider:${id}`);
  let provider = globalData.get(id);
  return provider ?? null;
}

export async function updateData(id, updates) {
  await fakeNetwork();
  let provider = globalData.get(id);
  if (!provider) throw new Error("No contact found for", id);
  globalData.set(id, {
    name: provider.name,
    description: provider.description,
    data: {
      "2d": [[...provider.data["2d"], [Math.random(), Math.random()]]],
      "3d": [
        [...provider.data["2d"], [Math.random(), Math.random(), Math.random()]],
      ],
    },
  });
  return new Map(globalData);
}

// fake a cache so we don't slow down stuff we've already seen
let fakeCache = {};

async function fakeNetwork(key) {
  if (!key) {
    fakeCache = {};
  }

  if (fakeCache[key]) {
    return;
  }

  fakeCache[key] = true;
  return new Promise((res) => {
    setTimeout(res, Math.random() * 800);
  });
}
