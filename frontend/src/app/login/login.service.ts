import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';
import { catchError, tap } from 'rxjs/operators';
import { throwError, BehaviorSubject } from 'rxjs';

import { User } from './user.model';

export interface AuthResponseData {
  token: string;
  username: string;
  expiresIn: string;
  admin: boolean;
}

@Injectable({ providedIn: 'root' })
export class LoginService {

  user = new BehaviorSubject<User>(null);
  private tokenExpirationTimer: any;

  constructor(private http: HttpClient, private router: Router) {}

  login(username: string, password: string) {
    return this.http
      .post<AuthResponseData>(
        'http://localhost:8080/login',
        {
          username: username,
          password: password,
        }
      )
      .pipe(
        catchError(this.handleError),
        tap(resData => {
          this.handleAuthentication(
            resData.username,
            resData.token,
            resData.admin,
            +resData.expiresIn
          );
        })
      );
  }

  private handleAuthentication(
    username: string,
    token: string,
    admin: boolean,
    expiresIn: number
  ) {
    const expirationDate = new Date(new Date().getTime() + expiresIn * 1000);
   // const user = new User(email, userId, token, expirationDate);

    const user = new User(username, admin, token, expirationDate);

    this.user.next(user);
    console.log('User to Local Storage: ', user);
   // this.autoLogout(expiresIn * 1000);
    localStorage.setItem('userData', JSON.stringify(user));
  }

  private handleError(errorRes: HttpErrorResponse) {
    let errorMessage = 'An unknown error occurred!';
    console.log('erro: ',errorRes.error.error);
    if (!errorRes.error || !errorRes.error.error) {
      return throwError(errorMessage);
    }
    if (errorRes.error.error == 'INVALID_PASSWORD'){
      errorMessage = 'This password is not correct.';
    }
    else if (errorRes.error.error == 'INVALID_JSON'){
      errorMessage = 'This password is not correct.';
    }
     else{
      errorMessage ="Banido durante " + errorRes.error.error + " segundos";

    }
    return throwError(errorMessage);
  }


}
