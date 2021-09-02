import { Injectable } from '@angular/core';
import { Subject } from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class AuthLibService {

  private _userName = new Subject<string>()

  constructor() { }

  get userName(): Subject<string> {
    return this._userName;
  }

  login(userName: string) {
    // console.log('Login: ', userName);
    this._userName.next(userName);
  }

}
