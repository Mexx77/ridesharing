describe('Account actions', () => {
    it('Logs in', () => {
        cy.visit('/');
        cy.get('[data-cy=enroll-btn]').click();
        cy.get('[data-cy=username-phone]').type('test');
        cy.get('[data-cy=password]').type('123456aB');
        cy.get('[data-cy=login-btn]').click();
    })
});
