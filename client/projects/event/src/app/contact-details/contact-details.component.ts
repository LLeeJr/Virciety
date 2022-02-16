import {Component, Inject, OnInit} from '@angular/core';
import {FormControl, Validators} from "@angular/forms";
import {emptyTextValidator} from "../create-event/create-event.component";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";

@Component({
  selector: 'app-contact-details',
  templateUrl: './contact-details.component.html',
  styleUrls: ['./contact-details.component.scss']
})
export class ContactDetailsComponent implements OnInit {
  firstname: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  lastname: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  street: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  housenumber: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  postalcode: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  city: FormControl = new FormControl('', [Validators.required, emptyTextValidator('')]);
  email: FormControl = new FormControl('', [Validators.required, Validators.email, emptyTextValidator('')]);

  constructor(@Inject(MAT_DIALOG_DATA) public data: string,
              private dialogRef: MatDialogRef<ContactDetailsComponent>) { }

  ngOnInit(): void {
  }


  fieldsValid(): boolean {
    return this.firstname.valid && this.lastname.valid && this.street.valid && this.housenumber.valid && this.postalcode.valid && this.city.valid && this.email.valid;
  }

  submit() {
    this.dialogRef.close({
      username: this.data,
      firstname: this.firstname.value as string,
      lastname: this.lastname.value as string,
      address: {
        street: this.street.value as string,
        housenumber: this.housenumber.value as string,
        postalcode: this.postalcode.value as string,
        city: this.city.value as string,
      },
      email: this.email.value as string,
    })
  }
}
