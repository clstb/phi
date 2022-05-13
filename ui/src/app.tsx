import React from "react";
import './styles/styles.css'
import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import {Classify, Home, LoginPage} from "./components";


// @ts-ignore
export const PrivateRoute = ({children}) => {
  if (sessionStorage.getItem("sessId")) {
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
                   <Home/>
                 </PrivateRoute>
               }
        />
        <Route path={'/login'}
               element={
                 <LoginPage/>
               }
        />
        <Route path={'/home'}
               element={
                 <PrivateRoute>
                   <Home/>
                 </PrivateRoute>
               }
        />
        <Route path={'/classify'}
               element={
                 <PrivateRoute>
                   <Classify/>
                 </PrivateRoute>
               }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
