import {Component, OnInit, ViewChild} from '@angular/core';
import {MatPaginator} from '@angular/material/paginator';
import {ManageService, ReadResponseData} from './manage.service';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import {LoginService, AuthResponseData} from '../login/login.service';
import { MatTableDataSource } from '@angular/material/table';
import {MatSort} from '@angular/material/sort';
import { MatTable } from '@angular/material/table';
import { MatDialog } from '@angular/material/dialog';
import { User } from '../login/user.model';

import { DialogBoxComponent } from '../dialog-box/dialog-box.component';


@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.css']
})
export class ManageComponent implements OnInit {
  error: string = null;
  displayedColumns: string[] = ['id', 'username', 'birth', 'name', 'family', 'admin', 'action'];
  dataSource ;
  isLoading = false;
  isAdmin;
  private userSub: User;


  // tslint:disable-next-line: max-line-length
  constructor(private manageService: ManageService, private router: Router , private loginService: LoginService, public dialog: MatDialog )  {

    }

    @ViewChild(MatTable, {static: true}) table: MatTable<any>;
    @ViewChild(MatPaginator, {static: true}) paginator: MatPaginator;
    @ViewChild(MatSort, {static: true}) sort: MatSort;


  ngOnInit() {
    this.userSub = JSON.parse(localStorage.getItem('userData'));
    this.isAdmin = this.userSub.Admin;
    console.log(this.isAdmin);
    this.readUsers();

  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }


readUsers() {
  let readObs: Observable<ReadResponseData[]>;
  this.isLoading = true;
  readObs = this.manageService.readAllUsers();
  readObs.subscribe(
    resData => {
      this.isLoading = false;
      console.log("resData ", resData);
      this.dataSource = new MatTableDataSource(resData);
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;

     // return resData;
      //this.dataSource = resData;
     // console.log("resposta: ", resData[0].username);
    },
    errorMessage => {
      this.error = errorMessage;
      //return null;
    }
  );
}

insertUsers(usernameF: string, passwordF: string, birthF: string, nameF: string , familyF: string, roleF: string){
  let readObs: Observable<any>;
  this.isLoading = true;
  if (!roleF){
    roleF = "normal";
  }
  readObs = this.manageService.newUser(usernameF, passwordF, birthF, nameF, familyF, roleF);
  readObs.subscribe(
    resData => {
      this.isLoading = false;
      console.log(resData);
    },
    errorMessage => {
      this.error = errorMessage;
            //return null;
    }
  );

}

updateUsers(id :number,usernameF: string, birthF: string, nameF: string , familyF: string, roleF: string){
  let readObs: Observable<any>;
  this.isLoading = true;
  if (!roleF){
    roleF = "normal";
  }
  readObs = this.manageService.updateUser(id,usernameF,  birthF, nameF, familyF, roleF);
  readObs.subscribe(
    resData => {
      this.isLoading = false;
      console.log(resData);
    },
    errorMessage => {
      this.error = errorMessage;
    }
  );
  }
  deleteUser(usernameF: string){
    let readObs: Observable<any>;
    this.isLoading = true;
    readObs = this.manageService.deleteUser(usernameF);
    readObs.subscribe(
      resData => {
        this.isLoading = false;
        console.log(resData);
      },
      errorMessage => {
        this.error = errorMessage;
      }
    );
    }
openDialog(action,obj) {
  obj.action = action;
  const dialogRef = this.dialog.open(DialogBoxComponent, {
    width: '250px',
    height: '550px',
    data:obj
  });
  dialogRef.afterClosed().subscribe(result => {
    if(result.event == 'Add'){
      // tslint:disable-next-line: max-line-length
      console.log('test', result.data);
      // tslint:disable-next-line: max-line-length
      this.insertUsers(result.data.username, result.data.password, result.data.birth, result.data.name, result.data.family, result.data.admin);
      this.readUsers();

      }
      else if(result.event == 'Update'){
        console.log('test', result.data);
        // tslint:disable-next-line: max-line-length
        this.updateUsers(result.data.id,result.data.username,result.data.birth, result.data.name, result.data.family, result.data.admin);
        this.readUsers();
      }else if(result.event == 'Delete'){
        console.log('test', result.data);
        // tslint:disable-next-line: max-line-length
        this.deleteUser(result.data.username);
        this.readUsers();    }
  });
}
addRowData(row_obj){
  console.log('new user: ', row_obj);
}
updateRowData(row_obj){
  console.log('new user: ', row_obj);
}
deleteRowData(row_obj){
  console.log('delete user: ', row_obj);
}
}
