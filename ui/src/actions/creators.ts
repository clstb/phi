import {
  ActionTypeEnum,
  IGetSessionAction,
  ILinkTinkAction,
  ILoginAction,
  IRegisterAction,
  ISetUsernameAction
} from "./types";
import {Dispatch} from "redux";


export const phiLogin = (id: string) => {
  return (dispatch: Dispatch<ILoginAction>) => {
    dispatch({
      type: ActionTypeEnum.LOGIN,
      sessionId: id
    })
  }
}

export const phiRegister = (id: string) => {
  return (dispatch: Dispatch<IRegisterAction>) => {
    dispatch({
      type: ActionTypeEnum.REGISTER,
      sessionId: id
    })
  }
}

export const saveUsername = (username: string) => {
  return (dispatch: Dispatch<ISetUsernameAction>) => {
    dispatch({
      type: ActionTypeEnum.SET_USERNAME,
      username: username
    })
  }
}


export const linkTink = (linked: boolean) => {
  return (dispatch: Dispatch<ILinkTinkAction>) => {
    dispatch({
      type: ActionTypeEnum.LINK_TINK,
      accountLinked: true
    })
  }
}

export const getSession = (id: string) => {
  return (dispatch: Dispatch<IGetSessionAction>) => {
    dispatch({
        type: ActionTypeEnum.GET_SESSION,
        sessionId: id
      }
    )
  }
}



