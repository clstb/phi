import React from 'react';
import ReactDOM from 'react-dom/client';
import {Provider} from 'react-redux'
import {store} from './store/store'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import {Classify, LoginPage, Home} from "./components";
import {CssBaseline} from "@mui/material";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);


root.render(
  <React.StrictMode>
    <Provider store={store}>
      <CssBaseline>
      <BrowserRouter>
        <Routes>
          <Route path={'/'} element={<LoginPage/>}
          />
          <Route path={'/home'} element={<Home/>}
          />
          <Route path={'/classify'} element={<Classify/>}
          />
        </Routes>
      </BrowserRouter>
      </CssBaseline>
    </Provider>
  </React.StrictMode>
);

