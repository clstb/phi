import {CORE_URI, DEFAULT_HEADERS} from "../constants";
import axios from "axios";
import {Button, FilledInput, FormGroup} from "@mui/material";
import React, {useContext, useState} from "react";
import '../styles/styles.css'
import {AppContext} from "../index";
import {NavigateFunction, useNavigate} from "react-router-dom";


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

const doLoginRegister = (path: string,
                         username: string,
                         password: string,
                         errSetter: ((err: any) => void),
                         navigate: NavigateFunction) => {
  axios.post(
    `${CORE_URI}/${path}`,
    {
      headers: DEFAULT_HEADERS,
      username: username,
      password: password
    })
    .then(
      res => {
        console.log(res.data)
        const sessId = res.data['sessionId']
        sessionStorage.setItem('username', username)
        sessionStorage.setItem('sessId', sessId)
        navigate('/home', {replace: true})
      }
    )
    .catch(err => {
      errSetter(err.response.statusText)
    })
}


export function LoginPage() {

  const [error, setError] = useState(undefined);
  if (error) {
    throw Error(error)
  }


  const context = useContext(AppContext)

  const {value: username, bind: bindUsername, reset: resetUsername} = useInput('');
  const {value: password, bind: bindPassword, reset: resetPassword} = useInput('');

  let navigate = useNavigate();

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
            <Button type={"submit"} onClick={(e) => {
              e.preventDefault()
              doLoginRegister("login", username, password, setError, navigate)
            }
            }>
              Login
            </Button>
            <Button type={"submit"} onClick={() => doLoginRegister("register", username, password, setError, navigate)
            }>
              Register
            </Button>
          </div>
        </div>
      </FormGroup>
    </div>
  );
}

