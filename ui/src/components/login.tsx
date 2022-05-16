import {CORE_URI, DEFAULT_HEADERS, LOGIN_PATH, REGISTER_PATH, SESS_ID, USERNAME} from "../constants";
import axios from "axios";
import {Button, FilledInput, FormGroup} from "@mui/material";
import React, {useState} from "react";
import '../styles/styles.css'
import {useNavigate} from "react-router-dom";
import {useInput} from "./util";



export function LoginPage() {

  const [error, setError] = useState<string>();
  if (error) {
    throw Error(error)
  }

  const {value: username, bind: bindUsername} = useInput();
  const {value: password, bind: bindPassword} = useInput();

  let navigate = useNavigate();


  const doLoginRegister = (path: string) => {
    if (!username || ! password){
      return
    }
    axios.post(
      `${CORE_URI}${path}`,
      {
        headers: DEFAULT_HEADERS,
        username: username,
        password: password
      })
      .then(
        res => {
          console.log(res.data)
          const sessId = res.data['sessionId']
          sessionStorage.setItem(USERNAME, username)
          sessionStorage.setItem(SESS_ID, sessId)
          navigate('/token')
        }
      )
      .catch(err => {
        setError(err.response.statusText)
      })
  }

  const doLogin = () => doLoginRegister(LOGIN_PATH)
  const doRegister = () => doLoginRegister(REGISTER_PATH)


  return (
    <div>
      <FormGroup>
        <div className='App-body'>
          <div className="input-container">
            <FilledInput type="text"
                         required={true}
                         placeholder={"username"}
                         {...bindUsername}
            />
          </div>
          <div className="input-container">
            <FilledInput type="password"
                         placeholder={"password"}
                         required={true}
                         {...bindPassword}
            />
          </div>
          <div className="button-container">
            <Button type={"submit"} onClick={doLogin}>
              Login
            </Button>
            <Button type={"submit"} onClick={doRegister}>
              Register
            </Button>
          </div>
        </div>
      </FormGroup>
    </div>
  );
}

