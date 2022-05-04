import {createStore, applyMiddleware} from "redux";
import reducers from "../reducers"
import thunk from 'redux-thunk'
import {middleware} from "./middleware";


export const store = createStore(
  reducers,
  {},
  applyMiddleware(thunk, middleware)
)

