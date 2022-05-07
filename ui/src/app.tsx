import React, {useContext} from "react";
import './styles/styles.css'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import {AppContext, PrivateRoute} from "./index";
import {Classify, Home, LoginPage} from "./components";


function App() {

  const value = useContext(AppContext)

  return (
    <BrowserRouter>
      <Routes>
        <Route path={'/'}
               element={value.sessionId? <Home/> : <LoginPage/>}
        />
        <Route path={'/login'}
               element={<LoginPage/>}
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
