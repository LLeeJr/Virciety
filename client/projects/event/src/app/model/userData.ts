export class UserData {
  private _username: string;
  private _firstname: string;
  private _lastname: string;
  private _address: Address;
  private _email: string;

  constructor(data: any) {
    this._username = data.username;
    this._email = data.email;
    this._firstname = data.firstname;
    this._lastname = data.lastname;
    this._address = {
      street: data.street,
      housenumber: data.housenumber,
      postalcode: data.postalcode,
      city: data.city,
    }
  }

  get username(): string {
    return this._username;
  }

  get firstname(): string {
    return this._firstname;
  }

  set firstname(value: string) {
    this._firstname = value;
  }

  get lastname(): string {
    return this._lastname;
  }

  set username(value: string) {
    this._username = value;
  }

  set lastname(value: string) {
    this._lastname = value;
  }

  get address(): Address {
    return this._address;
  }

  set address(value: Address) {
    this._address = value;
  }

  get email(): string {
    return this._email;
  }

  set email(value: string) {
    this._email = value;
  }
}

interface Address {
  street: string;
  housenumber: string;
  postalcode: string;
  city: string;
}
