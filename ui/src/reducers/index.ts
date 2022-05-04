import {combineReducers} from "redux";
import {reducer} from "./reducer";

const reducers = combineReducers({
  login: reducer
})

export default reducers

export * from "./state"
export * from "./reducer"
