import { matchSorter } from "match-sorter";
import sortBy from "sort-by";

export async function getProviders(query) {
  const resp = await fetch("http://localhost:5050/api/v1/providers");
  const respData = await resp.json();
  let providers = respData.providers;
  if (!providers) return [];
  if (query) {
    providers = matchSorter(providers, query, { keys: ["name"] });
  }
  return providers.sort(sortBy("name"));
}

export async function getProvider(uid) {
  const resp = await fetch("http://localhost:5050/api/v1/providers/" + uid);
  const provider = await resp.json();
  return provider ?? null;
}

export async function getProviderProjections(uid) {
  const resp = await fetch(
    "http://localhost:5050/api/v1/providers/" + uid + "/projections",
  );
  const embeddings = await resp.json();
  return embeddings ?? null;
}

export async function embedData(uid, updates) {
  await fetch("http://localhost:5050/api/v1/providers/" + uid + "/embeddings", {
    method: "PUT",
    body: JSON.stringify(updates),
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export async function deleteData(uid) {
  await fetch("http://localhost:5050/api/v1/providers/" + uid + "/embeddings", {
    method: "DELETE",
  });
}

export async function computeData(uid, updates) {
  await fetch(
    "http://localhost:5050/api/v1/providers/" + uid + "/projections",
    {
      method: "PATCH",
      body: JSON.stringify(updates),
      headers: {
        "Content-Type": "application/json",
      },
    },
  );
}
