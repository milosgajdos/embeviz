import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import Root, { loader as rootLoader } from "./routes/root";
import Embed, {
  loader as embedLoader,
  action as embedAction,
} from "./routes/embed";
import Index from "./routes/index";
import ErrorPage from "./error-page";

import "./index.css";

const router = createBrowserRouter(
  [
    {
      path: "/",
      element: <Root />,
      errorElement: <ErrorPage />,
      loader: rootLoader,
      children: [
        {
          errorElement: <ErrorPage />,
          children: [
            { index: true, element: <Index /> },
            {
              path: "provider/:uid",
              element: <Embed />,
              loader: embedLoader,
            },
            {
              path: "provider/:uid/update",
              element: <Embed />,
              action: embedAction,
              errorElement: <div>Oops! There was an error.</div>,
            },
          ],
        },
      ],
    },
  ],
  {
    basename: "/ui",
  }
);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
