import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders, HttpParams } from '@angular/common/http';
import { Router } from '@angular/router';
import { map, catchError, tap } from 'rxjs/operators';
import { throwError, BehaviorSubject } from 'rxjs';
import { User } from '../login/user.model';

import { Person } from './person.model';

export interface ReadResponseData {
  id: number;
  username: string;
  birth: string;
  name: string;
  family: string;
  admin: string;
}


@Injectable({ providedIn: 'root' })
export class ManageService {

  private userSub: User;

  person = new BehaviorSubject<Person>(null);


constructor(private http: HttpClient, private router: Router) {}

readAllUsers(){
  this.userSub = JSON.parse(localStorage.getItem('userData'));
  const token = this.userSub._token;
  const headers = new HttpHeaders({
    'Content-Type': 'application/json',
    Authorization: token });
  const options = { headers };
  return this.http
      .get<ReadResponseData[]>(
        'http://localhost:8080/read', options)
      .pipe(
        map(responseData => {
          return responseData;
        }),
        catchError(this.handleError)
      );
  }



  newUser(usernameF: string, passwordF: string, birthF: string, nameF: string , familyF: string, roleF: string ){

    this.userSub = JSON.parse(localStorage.getItem('userData'));

    console.log(this.userSub)
    const token = this.userSub._token;
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: token });
    const options = { headers };
    return this.http
        .put<ReadResponseData[]>(
          'http://localhost:8080/create',        {
            username: usernameF,
            password: passwordF,
            birth: birthF,
            name: nameF,
            family: familyF,
            role: roleF,
          } ,options)
        .pipe(
          map(responseData => {
            return responseData;
          }),
          catchError(this.handleError)
        );

  }


  updateUser(idf: number,usernameF: string,birthF: string, nameF: string , familyF: string, roleF: string ){

    this.userSub = JSON.parse(localStorage.getItem('userData'));

    console.log(this.userSub)
    const token = this.userSub._token;
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: token });
    const options = { headers };
    return this.http
        .post<ReadResponseData[]>(
          'http://localhost:8080/update',{
            id: idf,
            username: usernameF,
            birth: birthF,
            name: nameF,
            family: familyF,
            role: roleF,
          } ,options)
        .pipe(
          map(responseData => {
            return responseData;
          }),
          catchError(this.handleError)
        );

  }


    deleteUser(usernameF: string){

      this.userSub = JSON.parse(localStorage.getItem('userData'));

      console.log(this.userSub)
      const token = this.userSub._token;
      const headers = new HttpHeaders({
        'Content-Type': 'application/json',
        Authorization: token });
      const httpParams = new HttpParams().set('user', usernameF);
      const options = { headers,  params: httpParams };
      return this.http
          .delete<any>(
            'http://localhost:8080/delete'
           , options)
          .pipe(
            map(responseData => {
              return responseData;
            }),
            catchError(this.handleError)
          );


      }

  private handleError(errorRes: HttpErrorResponse) {
    let errorMessage = 'An unknown error occurred!';
    console.log('erro: ',errorRes.error);
    if (!errorRes.error || !errorRes.error.error) {
      return throwError(errorMessage);
    }
    if (errorRes.error == 'INVALID_TOKEN'){
      errorMessage = 'Token is not valid please relog';
    }
    else if (errorRes.error == 'INVALID_JSON'){
      errorMessage = 'Incorrect data parsed';
    }
    else if (errorRes.error == 'NO_PERMISSION'){
      errorMessage = 'No permission to handle this task';
    }
    else if (errorRes.error == 'INSERT_ERROR'){
      errorMessage = 'User could not be inserted';
    }
    else if (errorRes.error == 'UPDATE_ERROR'){
      errorMessage = 'Error updating user';
    }
    else if (errorRes.error == 'INVALID_USER'){
      errorMessage = 'User is not valid';
    }
    else if (errorRes.error == 'DELETE_ERROR'){
      errorMessage = 'Error on deleting';
    }
    return throwError(errorMessage);
  }


}
