import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import {APOLLO_NAMED_OPTIONS, NamedOptions} from "apollo-angular";
import {InMemoryCache} from "@apollo/client/core";
import {HttpClientModule} from "@angular/common/http";
import { PostCommentComponent } from './post-comment/post-comment.component';
import {HttpLink} from "apollo-angular/http";
import {FormsModule} from "@angular/forms";

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    PostCommentComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule
  ],
  providers: [
    {
      provide: APOLLO_NAMED_OPTIONS, // <-- Different from standard initialization
      useFactory(httpLink: HttpLink): NamedOptions {
        return {
          comment: {
            cache: new InMemoryCache(),
            link: httpLink.create({
              uri: 'http://localhost:8084/query',
            }),
          },
        };
      },
      deps: [HttpLink],
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
