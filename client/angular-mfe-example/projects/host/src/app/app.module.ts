import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import {Apollo} from "apollo-angular";
import {InMemoryCache} from "@apollo/client/core";
import {HttpClientModule} from "@angular/common/http";
import { PostCommentComponent } from './post-comment/post-comment.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    PostCommentComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor(
    apollo: Apollo
  ) {

    apollo.create({
      cache: new InMemoryCache(),
    });
  }
}
