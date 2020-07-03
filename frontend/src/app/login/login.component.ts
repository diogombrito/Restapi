import { Component } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

import {LoginService, AuthResponseData} from './login.service';
import { User } from './user.model';

@Component({
  selector: 'app-logincomponent',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent{
  isLoginMode = true;
  isLoading = false;
  error: string = null;
  checkoutForm;
  private userSub: User;

  constructor(private authService: LoginService, private router: Router, private formBuilder: FormBuilder    )  {
    this.checkoutForm = this.formBuilder.group({
      username: '',
      password: ''
    });

  }


  onSubmit(form) {
    if (!form.valid) {
      return;
    }
    const username = form.value.username;
    const password = form.value.password;

    let authObs: Observable<AuthResponseData>;

    this.isLoading = true;


    authObs = this.authService.login(username, password);

    console.log(username, password);
    authObs.subscribe(
      resData => {
        console.log("resposta: ", resData);
        this.isLoginMode = false;
        this.isLoading = false;
        this.router.navigate(['/manage']);
      },
      errorMessage => {
        console.log("error: ",errorMessage);
        this.error = errorMessage;
        this.isLoading = false;
      }
    );

    form.reset();
  }
}
