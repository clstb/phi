import {bindActionCreators} from 'redux';
import {actionCreators} from '../actions'
import {useDispatch, useSelector} from "react-redux";
import {CORE_URI, DEFAULT_HEADERS} from "../constants";
import axios from "axios";
import {IState} from "../reducers";
import {Button, FilledInput, FormGroup} from "@mui/material";
import '../styles/styles.css'
import Home from "./home";
import {useState} from "react";


export function LoginPage() {

  const useInput = (initialValue: string) => {
    const [value, setValue] = useState(initialValue);
    return {
      value,
      setValue,
      reset: () => setValue(""),
      bind: {
        value,
        onChange: (event: React.ChangeEvent<HTMLInputElement>) => {
          setValue(event.target.value);
        }
      }
    };
  };

  const {value: username, bind: bindUsername, reset: resetUsername} = useInput('');
  const {value: password, bind: bindPassword, reset: resetPassword} = useInput('');

  // @ts-ignore
  const state: IState = useSelector(state1 => state1['login'])
  console.log(state)
  const dispatch = useDispatch();
  const {phiLogin} = bindActionCreators(actionCreators, dispatch)


  const login = (e: React.MouseEvent) => {
    if (username.length === 0 || password.length === 0) {
      return
    }
    axios.post(
      `${CORE_URI}/login`,
      {
        headers: DEFAULT_HEADERS,
        username: username,
        password: password
      })
      .then(
        res => {
          console.log(res.data)
          var id = res.data['sessionId']
          phiLogin(username, id)
          resetUsername()
          resetPassword()
        }
      )
      .catch(alert)
    e.preventDefault()
  }

  const register = (e: React.MouseEvent) => {
    if (username.length === 0 || password.length === 0) {
      return
    }
    axios.post(
      `${CORE_URI}/register`,
      {
        headers: DEFAULT_HEADERS,
        username: username,
        password: password
      })
      .then(
        res => {
          console.log(res.data)
          var id = res.data['sessionId']
          phiLogin(username, id)
          resetUsername()
          resetPassword()
        }
      )
      .catch(alert)
    e.preventDefault()
  }

  if (state.sessionId) {
    return <Home/>
  }
  return (
    <div>
      <FormGroup>
        <div className='App-body'>
          <div className="input-container">
            <FilledInput type="text"
                         {...bindUsername}
                         placeholder={"username"}
                         required
            />
          </div>
          <div className="input-container">
            <FilledInput type="password"
                         placeholder={"password"}
                         {...bindPassword}
                         required
            />
          </div>
          <div className="button-container">
            <Button type={"submit"} onClick={login}>
              Login
            </Button>
            <Button type={"submit"} onClick={register}>
              Register
            </Button>
          </div>
        </div>
      </FormGroup>
    </div>
  );
}

