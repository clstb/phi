import {
  ActionTypeEnum,
  ILoginAction,
} from "./types";
import {Dispatch} from "redux";


export const phiLogin = (username: string, id: string) => {
  return (dispatch: Dispatch<ILoginAction>) => {
    dispatch({
      type: ActionTypeEnum.LOGIN,
      sessionId: id,
      username: username
    })
  }
}



