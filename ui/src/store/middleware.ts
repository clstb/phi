import { Middleware, Dispatch } from 'redux';
import {IAction} from "../actions/types";


export const middleware: Middleware = api => (next: Dispatch<IAction>) => action  => {
  return next(action);
};
