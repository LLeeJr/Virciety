import { Injectable } from '@angular/core';
import {MfeOptions} from "./mfe";

@Injectable({
  providedIn: 'root'
})
export class LookupService {
  lookup(): Promise<MfeOptions[]> {
    return Promise.resolve([
      {
        remoteEntry: 'http://localhost:5002/remoteEntry.js',
        remoteName: 'post',
        exposedModule: './CreatePostComponent',

        displayName: 'CreatePost',
        componentName: 'CreatePostComponent',
      },
    ] as MfeOptions[])
  }
}
