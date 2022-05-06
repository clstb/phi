import {Reducer} from "redux";
import {ActionType, ActionTypeEnum} from "../actions/types";
import {initState, IState} from "./state";


export const reducer: Reducer<IState, ActionType> = (state: IState = initState, action: any) => {
  console.log(`Reducer received ${action.type}`)
  switch (action.type) {
    case ActionTypeEnum.LOGIN:
      return {
        ...state,
        sessionId: action.sessionId,
        username: action.username
      };
    default:
      return state;
  }
}









