import test from 'node:test';
import assert from 'node:assert/strict';
import {normalizeSpec, validateAnswers, restoreAnswers} from '../skills/grill-form/assets/model.mjs';

const input={title:'Review',questions:[{number:8,prompt:'Ship?',recommendation:'Review first.',options:['Yes','Later'],context:'  code\n'}]};

test('simple committed question format normalizes stable IDs and rich source',()=>{
  const spec=normalizeSpec(input,{richText:true});
  assert.equal(spec.questions[0].id,'q8');
  assert.equal(spec.questions[0].context,'  code\n');
  const result=validateAnswers(spec,[{questionId:'q8',choiceId:'option2',notes:'  next week  '}]);
  assert.deepEqual(result,[{questionId:'q8',number:8,question:'Ship?',choiceId:'option2',choice:'Later',notes:'next week'}]);
});

test('notes-only answers work; empty answers and Other without notes cannot submit',()=>{
  const spec=normalizeSpec(input);
  for(const choiceId of [null,'other'])assert.throws(()=>validateAnswers(spec,[{questionId:'q8',choiceId,notes:' '}]),/Missing answer/);
  assert.equal(validateAnswers(spec,[{questionId:'q8',choiceId:null,notes:'A different plan'}])[0].choice,null);
  assert.equal(validateAnswers(spec,[{questionId:'q8',choiceId:'other',notes:'Explain'}])[0].choice,'Something else');
});

test('reopen preserves partial notes but refuses a mismatched questionnaire',()=>{
  const spec=normalizeSpec(input);
  const answers=[{questionId:'q8',choiceId:null,notes:'  unfinished  '}];
  const draft={schemaVersion:1,spec,answers};
  assert.deepEqual(restoreAnswers(spec,draft),answers);
  const changed=structuredClone(spec);changed.questions[0].options.reverse();
  assert.throws(()=>restoreAnswers(changed,draft),/different questionnaire/);
  assert.throws(()=>restoreAnswers(spec,{...draft,answers:[{...answers[0],choiceId:'invented'}]}),/Invalid choice/);
  assert.throws(()=>validateAnswers(spec,[]),/one answer/);
});

test('invalid spec and reserved IDs fail before a usable questionnaire is shown',()=>{
  assert.throws(()=>normalizeSpec({...input,questions:[input.questions[0],input.questions[0]]}),/unique/);
  assert.throws(()=>normalizeSpec({...input,questions:[{...input.questions[0],id:'other'}]}),/Invalid question id/);
  assert.throws(()=>normalizeSpec({...input,questions:[{...input.questions[0],options:['same','SAME']}]}),/Duplicate option/);
});
