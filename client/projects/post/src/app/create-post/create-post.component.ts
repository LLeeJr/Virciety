import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {GQLService} from "../service/gql.service";
import {AuthLibService} from "auth-lib";
import {DataLibService} from "data-lib";

@Component({
  selector: 'app-create-post',
  templateUrl: './create-post.component.html',
  styleUrls: ['./create-post.component.scss']
})
export class CreatePostComponent implements OnInit {
  fileBase64: any;
  file: any;
  fileBackend: any;
  description: string = '';
  content_type: string = '';

  constructor(private gqlService: GQLService,
              private authService: AuthLibService,
              private dataService: DataLibService,
              private cd: ChangeDetectorRef) {
  }

  ngOnInit(): void { }

  onFileSelected(event: any) {
    this.fileBackend = null;
    this.file = event.target.files[0];

    if (this.file) {
      const reader = new FileReader();
      reader.readAsDataURL(this.file);

      reader.onload = () => {
        if (reader.result) {
          const base64 = reader.result;
          const data: string = base64.toString().split(";base64,")[0];

          this.content_type = data.split(":")[1];

          this.fileBase64 = base64;
        }
      }
    }
  }

  createPost() {
    if (this.fileBase64) {
      this.gqlService.createPost(this.fileBase64, this.description, this.authService.userName);
    }
  }

  alertFunction() {
    alert(`Content-Type ${this.content_type} is not supported`);
    this.content_type = '';
    this.cd.detectChanges();
  }
}
