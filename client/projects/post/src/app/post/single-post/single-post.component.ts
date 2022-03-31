import {Component, OnInit} from '@angular/core';
import {Post} from "../../model/post";
import {GQLService} from "../../service/gql.service";
import {ActivatedRoute} from "@angular/router";
import {Location} from "@angular/common";
import {AuthLibService} from "auth-lib";
import {Comment} from "../../model/comment";
import {KeycloakService} from "keycloak-angular";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

@Component({
  selector: 'app-single-post',
  templateUrl: './single-post.component.html',
  styleUrls: ['./single-post.component.scss'],
  exportAs: 'SinglePostComponent'
})
export class SinglePostComponent implements OnInit {

  post: Post | null;
  valid: boolean = true;
  comment: string = '';
  nameSourceMap: Map<string, any> = new Map<string, any>();
  source: string = '';
  username: string = '';
  isPhonePortrait: boolean = false;

  constructor(private auth: AuthLibService,
              private gqlService: GQLService,
              private route: ActivatedRoute,
              private location: Location,
              private keycloak: KeycloakService,
              private responsive: BreakpointObserver) { }

  async ngOnInit(): Promise<void> {
    this.responsive.observe(Breakpoints.HandsetPortrait).subscribe((result) => {
      this.isPhonePortrait = result.matches;
    })

    let keycloakProfilePromise = await this.keycloak.loadUserProfile();
    this.username = <string>keycloakProfilePromise.username;

    let postID = this.route.snapshot.paramMap.get('id');
    // get postID when opened via dialog
    if (postID === null) {
      postID = this.location.path().substring(3);
    }

    if (postID !== null) {
      let returnedData = this.gqlService.getPostByID(postID);

      if (returnedData instanceof Post) {
        this.post = returnedData;
        this.getData(this.post);
      } else {
        returnedData.subscribe({
          next: ({data}: any) => {
            this.post = new Post(data.getPost);

            this.getData(this.post);
          },
          error: (_: any) => {
            this.valid = false;
          }
        });
      }
    } else {
      this.valid = false
    }
  }

  addComment() {
    if (this.post) {
      const addCommentRequest = {
        postID: this.post.id,
        comment: this.comment,
        createdBy: this.auth.userName,
      };

      this.gqlService.addComment(this.post, addCommentRequest).subscribe({
        next: ({data}: any) => {
          const comment = new Comment(data.addComment);

          if (!this.nameSourceMap.has(this.username)) {
            this.getCurrentUserProfilePicture();
          }

          if (this.post)
            this.post.comments = [comment, ...this.post.comments];

          this.comment = '';
          // console.log('AddCommentData: ', data)
        },
        error: (error: any) => {
          console.error('there was an error sending the addComment-mutation', error);
        }
      });
    }
  }

  private getData(post: Post) {
    this.gqlService.getPostComments(post);

    this.auth.getUserByName(post.username).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        let {profilePictureId} = value.data.getUserByName;
        if (profilePictureId && profilePictureId !== '') {
          this.auth.getProfilePicture(profilePictureId).subscribe(picture => {
            if (picture && picture.data && picture.data.getProfilePicture) {
              this.source = picture.data.getProfilePicture;
            }
          })
        }
      }
    });
    this.gqlService.userProfilePictureIds.subscribe(map => {
      for (let entry of map.entries()) {
        if (entry[1] !== "") {
          this.auth.getProfilePicture(entry[1]).subscribe(value => {
            if (value && value.data && value.data.getProfilePicture) {
              this.nameSourceMap.set(entry[0], value.data.getProfilePicture);
            }
          })
        } else {
          this.nameSourceMap.set(entry[0], "");
        }
      }
    });
    let find = post.comments.find(comment => comment.createdBy === this.username);
    if (find != undefined && !this.nameSourceMap.has(this.username)) {
      this.getCurrentUserProfilePicture();
    }
  }

  private getCurrentUserProfilePicture() {
    this.auth.getUserByName(this.username).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        let {profilePictureId} = value.data.getUserByName;
        if (profilePictureId && profilePictureId !== '') {
          this.auth.getProfilePicture(profilePictureId).subscribe(picture => {
            if (picture && picture.data && picture.data.getProfilePicture) {
              this.nameSourceMap.set(this.username, picture.data.getProfilePicture);
            }
          })
        }
      }
    });
  }
}
