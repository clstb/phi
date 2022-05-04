import './reducer'
import reducers from "./index";

export interface IState {
  sessionId? : string,
  username?: string,
  accountLinked: boolean
}

export const initState: IState = {
  accountLinked: false
}

export type State = ReturnType<typeof reducers>
