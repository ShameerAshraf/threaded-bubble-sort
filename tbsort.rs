// sort only positive numbers

use std::thread;
use std::time::Duration;
use std::sync::mpsc;
use std::convert::TryInto;

fn main() {
    let array: [i32; 10] = [6, 4, 9, 2, 8, 1, 7, 5, 3 , 10];
    tbsort(array);
}

fn tbsort(array: [i32; 10]) {
    
    for val in array {  // for each 2nd val

        thread::spawn(move || {
            thread::sleep(Duration::from_secs(val.try_into().unwrap()));
        });
    }


    println!("Final state: {:?}", vec);

}
// arc to share array
// spawn for every 2nd elem starting 0
// bound check for thread for n - 1 or n - 2, dont spawn for n - 1 if array size odd?
// loop n times
// Barrier to sync threads 
// OR
// atomic counter, all threads decrease, signal wake up when zero
