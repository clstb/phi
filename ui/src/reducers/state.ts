import './reducer'
import reducers from "./index";

export interface IState {
  sessionId? : string
  username?: string
}

export const initState: IState = {
  sessionId: sessionStorage.getItem("phiSessionId") || undefined,
  username: sessionStorage.getItem("phiUsername") || undefined
}

export type AppState = ReturnType<typeof reducers>
