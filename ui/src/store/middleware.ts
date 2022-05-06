import {Dispatch, Middleware, MiddlewareAPI} from 'redux';
import {ActionTypeEnum, IAction, ILoginAction} from "../actions/types";


export const middleware: Middleware = (api: MiddlewareAPI<any>) => (next: Dispatch<IAction>) => (action: IAction)  => {
  if(action.type === ActionTypeEnum.LOGIN){
    sessionStorage.setItem("phiSessionId", (action as ILoginAction).sessionId);
    sessionStorage.setItem("phiUsername", (action as ILoginAction).username)
  }
  return next(action)
}
