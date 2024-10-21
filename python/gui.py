import io
import tkinter as tk
from tkinter import messagebox, filedialog, ttk
import grpc
import pet_pb2
import pet_pb2_grpc
from PIL import Image, ImageTk


class PetApp:
    def __init__(self, root):
        self.image_data = None
        self.root = root
        self.root.title("Pet Registration and Search")

        self.create_widgets()

        # Initialize gRPC client
        self.channel = grpc.insecure_channel('localhost:50051')
        self.stub = pet_pb2_grpc.PetServiceStub(self.channel)

        # List to hold image references to prevent garbage collection
        self.image_refs = []

    def create_widgets(self):
        # Registration Frame
        self.registration_frame = tk.Frame(self.root)
        self.registration_frame.pack(pady=10)

        tk.Label(self.registration_frame, text="Pet Name:").grid(row=0, column=0)
        self.pet_name_entry = tk.Entry(self.registration_frame)
        self.pet_name_entry.grid(row=0, column=1)

        tk.Label(self.registration_frame, text="Gender:").grid(row=1, column=0)
        self.pet_gender_entry = ttk.Combobox(self.registration_frame, values=["Male", "Female", "Null"])
        self.pet_gender_entry.grid(row=1, column=1)
        self.pet_gender_entry.current(0)  # Set default value to "Male"

        tk.Label(self.registration_frame, text="Age:").grid(row=2, column=0)
        self.pet_age_entry = tk.Entry(self.registration_frame)
        self.pet_age_entry.grid(row=2, column=1)

        tk.Label(self.registration_frame, text="Breed:").grid(row=3, column=0)
        self.pet_breed_entry = tk.Entry(self.registration_frame)
        self.pet_breed_entry.grid(row=3, column=1)

        tk.Label(self.registration_frame, text="Picture:").grid(row=4, column=0)
        self.pet_picture_entry = tk.Entry(self.registration_frame)
        self.pet_picture_entry.grid(row=4, column=1)
        tk.Button(self.registration_frame, text="Browse", command=self.browse_picture).grid(row=4, column=2)

        tk.Button(self.registration_frame, text="Register Pet", command=self.register_pet).grid(row=5, column=0,
                                                                                                columnspan=3)

        # Search Frame
        self.search_frame = tk.Frame(self.root)
        self.search_frame.pack(pady=10)

        tk.Label(self.search_frame, text="Select Search Field:").grid(row=0, column=0)
        self.search_field = ttk.Combobox(self.search_frame, values=["Name", "Gender", "Age", "Breed"])
        self.search_field.grid(row=0, column=1)
        self.search_field.current(0)  # Set default value to first option

        tk.Label(self.search_frame, text="Search Term:").grid(row=1, column=0)
        self.search_entry = tk.Entry(self.search_frame)
        self.search_entry.grid(row=1, column=1)

        tk.Button(self.search_frame, text="Search", command=self.search_pet).grid(row=1, column=2)

        # Results Frame
        self.results_frame = tk.Frame(self.search_frame)
        self.results_frame.grid(row=2, column=0, columnspan=3)

    def browse_picture(self):
        file_types = [("Image Files", "*.jpg;*.jpeg;*.png"), ("All Files", "*.*")]
        file_path = filedialog.askopenfilename(filetypes=[('All Files', '*.*')])
        if file_path:
            self.pet_picture_entry.delete(0, tk.END)
            self.pet_picture_entry.insert(0, file_path)

            try:
                # Open the image file and load into self.image_data as binary
                image = Image.open(file_path)
                self.image_data = io.BytesIO()
                image.save(self.image_data, format='PNG')  # Save as PNG to self.image_data
                self.image_data.seek(0)  # Reset the pointer to the start of the stream
            except Exception as e:
                messagebox.showerror("Error", f"Failed to load image: {str(e)}")

    def register_pet(self):
        name = self.pet_name_entry.get()
        gender = self.pet_gender_entry.get()
        age = self.pet_age_entry.get()
        breed = self.pet_breed_entry.get()

        if not all([name, gender, age, breed, self.image_data]):
            messagebox.showwarning("Input Error", "Please fill out all fields and select a picture.")
            return

        # Reset image data stream position
        self.image_data.seek(0)

        # Create a Pet object
        pet = pet_pb2.RegisterNewPetRequest(
            name=name,
            gender=gender,
            age=int(age),
            breed=breed,
            picture=self.image_data.read()  # Read the binary image data directly
        )

        # Call gRPC register method
        try:
            response = self.stub.RegisterNewPet(pet)
            print(response)
            # messagebox.showinfo("Registration", response.message)
        except grpc.RpcError as e:
            messagebox.showerror("Registration Error", f"An error occurred during registration: {e.details()}")

    def search_pet(self):
        search_term = self.search_entry.get()
        search_field = self.search_field.get()

        if not search_term:
            messagebox.showwarning("Input Error", "Please enter a search term.")
            return
        else:
            # Construct search request based on the selected field
            search_request = None
            try:
                if search_field == "Name":
                    search_request = pet_pb2.SearchPetRequest(name=search_term)
                elif search_field == "Gender":
                    search_request = pet_pb2.SearchPetRequest(gender=search_term)
                elif search_field == "Age":
                    search_request = pet_pb2.SearchPetRequest(age=int(search_term))  # Ensure age is an integer
                elif search_field == "Breed":
                    search_request = pet_pb2.SearchPetRequest(breed=search_term)
            except ValueError:
                messagebox.showerror("Input Error", "Age must be an integer.")
                return

            # Call gRPC search method
            try:
                response = self.stub.SearchPet(search_request)
            except grpc.RpcError as e:
                messagebox.showerror("Search Error", f"An error occurred during search: {e.details()}")
                return

            # Clear previous results
            for widget in self.results_frame.winfo_children():
                widget.destroy()

            self.image_refs.clear()  # Clear previous image references

            if response.pets:
                row = 0
                for pet in response.pets:
                    info = f"Name: {pet.name}, Gender: {pet.gender}, Age: {pet.age}, Breed: {pet.breed}"
                    tk.Label(self.results_frame, text=info).grid(row=row, column=0, sticky='w')
                    # Load image from bytes
                    image_data = io.BytesIO(pet.picture)
                    try:
                        pil_image = Image.open(image_data)
                        pil_image = pil_image.resize((100, 100))  # Resize image for display
                        tk_image = ImageTk.PhotoImage(pil_image)
                        image_label = tk.Label(self.results_frame, image=tk_image)
                        image_label.grid(row=row, column=1)
                        self.image_refs.append(tk_image)  # Keep a reference to avoid garbage collection
                    except Exception as e:
                        messagebox.showerror("Image Error", f"Failed to load image for {pet.name}: {str(e)}")
                    row += 1
            else:
                tk.Label(self.results_frame, text="No pets found.").grid(row=0, column=0)


if __name__ == "__main__":
    root = tk.Tk()
    app = PetApp(root)
    root.mainloop()
