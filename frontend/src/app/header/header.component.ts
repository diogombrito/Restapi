import {ChangeDetectorRef, ChangeDetectionStrategy,  Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';

import { LoginService } from '../login/login.service';
import { User } from '../login/user.model';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush


})
export class HeaderComponent implements OnInit, OnDestroy {
  isAuthenticated = false;
  private userSub: User;
  //private userSub2: Subscription;
  numberOfTicks;
  constructor(
    private authService: LoginService,private ref: ChangeDetectorRef
  ) {            setInterval(() => {
    this.numberOfTicks++;
    // the following is required, otherwise the view will not be updated
    this.ref.markForCheck();
  }, 1000);
  }

  ngOnInit() {
    this.userSub = JSON.parse(localStorage.getItem('userData'));
    console.log('is logged: ', this.userSub._token );
    console.log('user token: ', this.userSub );
    this.isAuthenticated = (!!this.userSub);
    console.log('user token: ', this.isAuthenticated );
    this.ref.detectChanges();
  }
/*     this.userSub2 = this.authService.user.subscribe(user => {
      this.isAuthenticated = !!user;
      console.log(!!user);
      console.log("user: ",user);

    });
 */



  ngOnDestroy() {
     localStorage.clear();
     // this.userSub2.unsubscribe();

  }

}
