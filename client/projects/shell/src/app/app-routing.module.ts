import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {HomeComponent} from "./home/home.component";
import {AuthGuard} from "./guard/auth.guard";
import {loadRemoteModule} from "@angular-architects/module-federation";

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    pathMatch: 'full',
  },
  {
    path: 'chat',
    canActivate: [AuthGuard],
    loadChildren: () =>
      loadRemoteModule({
        remoteName: 'chat',
        exposedModule: './ChatModule',
      }).then(m => m.ChatModule),
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
