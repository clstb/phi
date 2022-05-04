import {Reducer} from "redux";
import {ActionType, ActionTypeEnum} from "../actions/types";
import {initState, IState} from "./state";


export const reducer: Reducer<IState, ActionType> = (state: IState = initState, action: any) => {
  console.log(`Reducer received ${action.type}`)
  switch (action.type) {
    case ActionTypeEnum.SET_USERNAME:
      return {
        ...state,
        username: action.username
      }
    case ActionTypeEnum.LOGIN:
      return {
        ...state,
        sessionId: action.sessionId,
        loggedId: true
      };
    case ActionTypeEnum.REGISTER:
      return {
        ...state,
        sessionId: action.sessionId,
        loggedId: true
      }
    case ActionTypeEnum.LINK_TINK:
      return {
        ...state,
        accountLinked: action.accountLinked
      }
    case ActionTypeEnum.GET_SESSION:
      return {
        ...state,
        sessionId: action.sessionId
      }
    default:
      return state;
  }
}









