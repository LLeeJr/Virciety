import {Component, OnDestroy, OnInit} from '@angular/core';
import { Post } from "../model/post";
import {GQLService} from "../service/gql.service";
import {Observable} from "rxjs";
import {KeycloakService} from "keycloak-angular";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit, OnDestroy {
  posts: Observable<Post[]> | undefined;
  username: string;
  filter: string | null;

  constructor(private gqlService: GQLService,
              private keycloak: KeycloakService) {
  }

  async ngOnInit(): Promise<void> {
    await this.keycloak.isLoggedIn().then(loggedIn => {
      if (loggedIn) {
        this.keycloak.loadUserProfile().then(() => {
          this.username = this.keycloak.getUsername();

          this.filter = null; // TODO change to get username out of url
          this.posts = this.gqlService.getPosts(this.filter);

          this.gqlService.getPostCreated();
        })
      } else {
        this.keycloak.login();
      }
    });

    /*this.posts = this.gqlService.getPosts();
    this.gqlService.getPostCreated();*/
  }

  ngOnDestroy(): void {
    GQLService._oldestPostReached = false;
  }

  get oldestPostReached(): boolean {
    return GQLService._oldestPostReached;
  }

  onScroll() {
    if (!this.oldestPostReached) {
      this.gqlService.refreshPosts();
    }
  }

  addComment(post: Post, newComment: string) {
    const addCommentRequest = {
      postID: post.id,
      comment: newComment,
      createdBy: this.username
    };

    this.gqlService.addComment(post, addCommentRequest);
  }

  triggerEvent(post: Post, event: string) {
    if (event === 'like') {
      this.likePost(post);
    } else if (event === 'edit') {
      this.editPost(post);
    } else if (event === 'remove') {
      this.removePost(post);
    } else {
      this.showComments(post);
    }
  }

  likePost(post: Post) {
    const username: string = this.username;
    let liked: boolean = true;

    // check if it's a like or unlike
    const index = post.likedBy.indexOf(username, 0);
    if (index > -1) {
      const newLikedBy: string[] = []
      post.likedBy.forEach((user, i) => {
        if (index !== i) {
          newLikedBy.push(user);
        }
      });
      post.likedBy = newLikedBy;
      liked = false;
    } else {
      const newLikedBy: string[] = [username];
      post.likedBy.forEach(user => newLikedBy.push(user));
      post.likedBy = newLikedBy;
    }

    this.gqlService.likePost(post, liked);
  }

  editPost(post: Post) {
    post.editMode = false;
    this.gqlService.editPost(post);
  }

  removePost(post: Post) {
    post.editMode = false;
    this.gqlService.removePost(post);
  }

  showComments(post: Post) {
    post.commentMode = true;
    this.gqlService.getPostComments(post);
  }
}
