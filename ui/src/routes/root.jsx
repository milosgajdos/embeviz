import {
  Outlet,
  NavLink,
  Form,
  useLoaderData,
  useNavigation,
  useSubmit,
} from "react-router-dom";
import { getProviders } from "../lib/embeddings";
import { useEffect } from "react";

export async function loader({ request }) {
  const url = new URL(request.url);
  const q = url.searchParams.get("q");
  try {
    const providers = await getProviders(q);
    return { providers, q };
  } catch (error) {
    throw new Response("", {
      status: error.Status,
      statusText: "Fetching providers failed!",
    });
  }
}

export default function Root() {
  const { providers, q } = useLoaderData();
  const navigation = useNavigation();
  const submit = useSubmit();

  const searching =
    navigation.location &&
    new URLSearchParams(navigation.location.search).has("q");

  useEffect(() => {
    document.getElementById("q").value = q;
  }, [q]);

  return (
    <>
      <div id="sidebar">
        <h1>Embeviz</h1>
        <div>
          <Form method="get" id="search-form" role="search">
            <input
              id="q"
              className={searching ? "loading" : ""}
              placeholder="Search"
              type="search"
              name="q"
              defaultValue={q}
              onChange={(event) => {
                const isFirstSearch = q == null;
                submit(event.currentTarget.form, {
                  replace: !isFirstSearch,
                });
              }}
            />
            <div id="search-spinner" aria-hidden hidden={!searching}></div>
            <div className="sr-only" aria-live="polite"></div>
          </Form>
        </div>
        <nav>
          {providers.length ? (
            <ul>
              {providers.map((provider) => (
                <li key={provider.id}>
                  <NavLink
                    to={`provider/${provider.id}`}
                    className={({ isActive, isPending }) =>
                      isActive ? "active" : isPending ? "pending" : ""
                    }
                  >
                    {provider.name ? <>{provider.name}</> : <i>No Name</i>}{" "}
                  </NavLink>
                </li>
              ))}
            </ul>
          ) : (
            <p>
              <i>No providers</i>
            </p>
          )}
        </nav>
      </div>
      <div
        id="detail"
        className={navigation.state === "loading" ? "loading" : ""}
      >
        <Outlet />
      </div>
    </>
  );
}
