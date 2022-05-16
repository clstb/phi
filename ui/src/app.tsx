import React from "react";
import './styles/styles.css'
import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import {LoginPage, TokenPage} from "./components";
import {SESS_ID} from "./constants";


// @ts-ignore
export const PrivateRoute = ({children}) => {
  if (sessionStorage.getItem(SESS_ID)) {
    return children
  }
  return <Navigate to="/login"/>
}

function App() {

  return (
    <BrowserRouter>
      <Routes>
        <Route path={'/'}
               element={
                 <PrivateRoute>
                   <TokenPage/>
                 </PrivateRoute>
               }
        />
        <Route path={'/token'}
               element={
                 <PrivateRoute>
                   <TokenPage/>
                 </PrivateRoute>
               }
        />
        <Route path={'/login'}
               element={
                 <LoginPage/>
               }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
