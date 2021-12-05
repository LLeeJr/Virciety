import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {HomeComponent} from "./home/home.component";
import {AuthGuard} from "./guard/auth.guard";
import {loadRemoteModule} from "@angular-architects/module-federation";

const routes: Routes = [
  {
    path: '',
    redirectTo: '/home',
    pathMatch: 'full'
  },
  {
    path: 'chat',
    canActivate: [AuthGuard],
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'chat',
        exposedModule: './ChatModule',
      }).then(m => m.ChatModule),
  },
  {
    path: 'home',
    canActivate: [AuthGuard],
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'post',
        exposedModule: './PostModule',
      }).then(m => m.PostModule),
  },
  {
    path: 'user',
    canActivate: [AuthGuard],
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'user',
        exposedModule: './UserModule',
      }).then(m => m.UserModule),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
