import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import { RegisterComponent } from './register.component';
import {AlertService} from "../../../_services/alert.service";
import {AuthenticationService} from "../../_services/authentication.service";
import {Router} from "@angular/router";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {RouterTestingModule} from "@angular/router/testing";
import {of, throwError} from "rxjs";
import {LoginComponent} from "../login/login.component";

describe('RegisterComponent', () => {
  let component: RegisterComponent;
  let fixture: ComponentFixture<RegisterComponent>;
  let alertService: AlertService;
  let authService: AuthenticationService;
  let router: Router;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RegisterComponent ],
      imports: [
        HttpClientTestingModule,
        RouterTestingModule,
      ]
    })
    .compileComponents();
  });

  it('should redirect if a user is logged in', () => {
    authService = TestBed.inject(AuthenticationService);
    router = TestBed.inject(Router);

    spyOnProperty(authService, 'currentUserValue').and.returnValue({});
    spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));

    fixture = TestBed.createComponent(RegisterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();

    expect(router.navigate).toHaveBeenCalled();
  });

  describe('when a user is not signed in', () => {
    beforeEach(() => {
      fixture = TestBed.createComponent(RegisterComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
      alertService = TestBed.inject(AlertService);
      authService = TestBed.inject(AuthenticationService);
      router = TestBed.inject(Router);
    });

    it('should create', () => {
      expect(component).toBeTruthy();
      expect(component.registerForm).not.toBeUndefined();
      expect(component.registerForm.value.firstName).not.toBeUndefined();
      expect(component.registerForm.value.lastName).not.toBeUndefined();
      expect(component.registerForm.value.email).not.toBeUndefined();
      expect(component.registerForm.value.password).not.toBeUndefined();
      expect(component.registerForm.value.confirmPassword).not.toBeUndefined();
    });

    describe('onRegister', () => {
      beforeEach(() => {
        component.registerForm.setValue({
          firstName: 'first',
          lastName: 'last',
          email: 'email@email.com',
          password: 'password',
          confirmPassword: 'password',
          username: 'username',
        });
      });

      it('should warn on invalid form data', () => {
        component.registerForm.reset();
        spyOn(alertService, 'warn');
        component.onRegister();
        expect(alertService.warn).toHaveBeenCalledOnceWith('Please fill out all forms correctly', {id: 'register-alert', autoClose: 'true'});
        expect(alertService.warn).toHaveBeenCalledTimes(1);
      });

      it('should require matching passwords', () => {
        spyOn(alertService, 'warn');
        component.registerForm.value.password += 'a';
        component.onRegister();
        expect(alertService.warn).toHaveBeenCalledOnceWith('Passwords do not match', {id: 'register-alert', autoClose: 'true'});
        expect(alertService.warn).toHaveBeenCalledTimes(1);
      });

      it('should work and show any errors if there were any', () => {
        spyOn(alertService, 'error');
        spyOn(authService, 'requestRegister').and.returnValue(of({status: 200, errors: [{text: 'text', info: 'info'}]}));
        component.onRegister();
        expect(authService.requestRegister).toHaveBeenCalledTimes(1);
        expect(alertService.error).toHaveBeenCalledTimes(1);
      });

      it('should work and show any server errors if there were any', () => {
        spyOn(alertService, 'error');
        spyOn(authService, 'requestRegister').and.returnValue(throwError({status: 400, errors: [{text: 'text', info: 'info'}]}));
        component.onRegister();
        expect(authService.requestRegister).toHaveBeenCalledTimes(1);
        expect(alertService.error).toHaveBeenCalledTimes(1);
      });

      it('should work and redirect on valid request', () => {
        spyOn(alertService, 'success');
        spyOn(router, 'navigate').and.returnValue(Promise.resolve(true));
        spyOn(authService, 'requestRegister').and.returnValue(of({status: 200}));
        component.onRegister();
        expect(authService.requestRegister).toHaveBeenCalledTimes(1);
        waitForAsync(() => expect(alertService.success).toHaveBeenCalledTimes(1));
      });
    });
  });
});
