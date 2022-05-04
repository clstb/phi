export enum ActionTypeEnum {
  LOGIN = 'LOGIN',
  REGISTER = 'REGISTER',
  LINK_TINK = 'LINK_TINK',
  SET_USERNAME = 'SET_USERNAME',
  GET_SESSION = 'GET_SESSION',
}

export interface IAction {
  type: ActionTypeEnum
}

export interface ISetUsernameAction {
  type: ActionTypeEnum.SET_USERNAME,
  username: string
}

export interface ILoginAction extends IAction {
  type: ActionTypeEnum.LOGIN,
  sessionId: string
}

export interface IGetSessionAction extends IAction {
  type: ActionTypeEnum.GET_SESSION,
  sessionId: string
}

export interface IRegisterAction extends IAction{
  type: ActionTypeEnum.REGISTER,
  sessionId: string
}

export interface ILinkTinkAction extends IAction{
  type: ActionTypeEnum.LINK_TINK
  accountLinked: boolean
}

export type ActionType = ILoginAction | IRegisterAction | ILinkTinkAction | IGetSessionAction | ISetUsernameAction
