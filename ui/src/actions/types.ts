export enum ActionTypeEnum {
  LOGIN = 'LOGIN'
}

export interface IAction {
  type: ActionTypeEnum
}

export interface ILoginAction extends IAction {
  type: ActionTypeEnum.LOGIN,
  sessionId: string,
  username: string
}

export type ActionType = ILoginAction;
